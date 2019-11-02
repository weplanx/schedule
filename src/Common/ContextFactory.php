<?php

namespace Lumen\Extra\Common;

use Lumen\Extra\Contracts\ContextInterface;

/**
 * Class ContextFactory
 * @package Lumen\Extra\Common
 */
final class ContextFactory implements ContextInterface
{
    /**
     * @var array
     */
    private $context = [];

    /**
     * @param string $abstract
     * @param mixed $value
     * @inheritDoc
     */
    public function set(string $abstract, $value)
    {
        $this->context[$abstract] = $value;
    }

    /**
     * @param $abstract
     * @return mixed
     * @inheritDoc
     */
    public function get($abstract)
    {
        return !empty($this->context[$abstract]) ?
            $this->context[$abstract] :
            null;
    }
}
