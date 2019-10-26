<?php

namespace lumen\extra\common;

use Illuminate\Support\Str;
use Lcobucci\JWT\Builder;
use Lcobucci\JWT\Parser;
use Lcobucci\JWT\Signer\Hmac\Sha256;
use Lcobucci\JWT\Signer\Key;
use Lcobucci\JWT\Token;
use lumen\extra\redis\RefreshToken;

final class TokenFactory
{
    /**
     * Token config
     * @var array $config
     */
    private $config;

    /**
     * App secret
     * @var string $secret
     */
    private $secret;

    /**
     * Token signer
     * @var Sha256 $signer
     */
    private $signer;

    /**
     * Token
     * @var Token
     */
    private $token;

    /**
     * JwtAuthFactory constructor.
     * @param string $secret App Key
     * @param array $config Jwt Config
     */
    public function __construct(string $secret,
                                array $config)
    {
        $this->secret = $secret;
        $this->config = $config;
        $this->signer = new Sha256();
    }

    /**
     * Get Token
     * @return Token
     */
    public function getToken(string $token = null)
    {
        return empty($token) ? $this->token : (new Parser())->parse($token);
    }

    /**
     * Set Token
     * @param string $scene Token scene
     * @param array $symbol Symbol Tag
     * @return bool|string
     * @throws \Exception
     */
    public function setToken(string $scene,
                             array $symbol = [])
    {
        if (empty($this->config[$scene])) {
            throw new \Exception('not exists scene: ' . $scene);
        }

        $jti = Str::uuid()->toString();
        $ack = Str::random();

        $this->token = (new Builder())
            ->issuedBy($this->config[$scene]['issuer'])
            ->permittedFor($this->config[$scene]['audience'])
            ->identifiedBy($jti, true)
            ->withClaim('ack', $ack)
            ->withClaim('symbol', $symbol)
            ->expiresAt(time() + $this->config[$scene]['expires'])
            ->getToken($this->signer, new Key($this->secret));

        if (!empty($this->config[$scene]['auto_refresh'])) {
            $result = (new RefreshToken)
                ->factory($jti, $ack, $this->config[$scene]['auto_refresh']);

            if ($result == false) {
                return false;
            }
        }

        return (string)$this->token;
    }

    /**
     * Token Verify
     * @param string $scene Token scene
     * @param string $token String Token
     * @return bool|string
     * @throws \Exception
     */
    public function verify(string $scene,
                           string $token)
    {
        if (empty($this->config[$scene])) {
            throw new \Exception('not exists scene: ' . $scene);
        }

        $this->token = (new Parser())->parse($token);

        if (!$this->token->verify($this->signer, $this->secret)) {
            return false;
        }

        if ($this->token->getClaim('iss') != $this->config[$scene]['issuer'] ||
            $this->token->getClaim('aud') != $this->config[$scene]['audience']) {
            return false;
        }

        if ($this->token->isExpired()) {
            if (empty($this->config[$scene]['auto_refresh'])) {
                return false;
            }

            $result = (new RefreshToken)->verify(
                $this->token->getClaim('jti'),
                $this->token->getClaim('ack')
            );

            if (!$result) {
                return false;
            }

            $newToken = (new Builder())
                ->issuedBy($this->config[$scene]['issuer'])
                ->permittedFor($this->config[$scene]['audience'])
                ->identifiedBy($this->token->getClaim('jti'), true)
                ->withClaim('ack', $this->token->getClaim('ack'))
                ->withClaim('symbol', $this->token->getClaim('symbol'))
                ->expiresAt(time() + $this->config[$scene]['expires'])
                ->getToken($this->signer, new Key($this->secret));

            $this->token = $newToken;
            return (string)$this->token;
        }

        return true;
    }
}
